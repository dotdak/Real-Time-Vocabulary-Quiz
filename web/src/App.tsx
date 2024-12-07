import React, { useEffect, useState } from 'react';
import './App.css';
import { ThemeProvider } from '@mui/material/styles';
import LoginBox from './components/LoginBox';
import { useParams, useNavigate } from 'react-router-dom';
import { theme } from './theme';
import { WebSocketReconnector } from './socket';
import {
  Button,
  Container,
  Divider,
  Grid,
  LinearProgress,
  Stack,
  Typography,
} from '@mui/material';
import LeaderBoard from './components/LeaderBoard';
import QuestionBoard from './components/QuestionBoard';
import client, { ApiClient } from './api-client';
function App() {
  const params = useParams();
  const quizId = params?.quizId || '';
  const navigate = useNavigate();
  const [username, setUsername] = useState(
    localStorage.getItem('username') || '',
  );

  const [openReconnecting, setOpenReconnecting] = useState<boolean>(false);
  const [socket, setSocket] = useState<WebSocketReconnector>();
  const [leaderBoard, setLeaderBoard] = useState({});
  const [boardStatus, setBoardStatus] = useState('inProgress');
  const [counter, setCounter] = useState(0);
  const [quiz, setQuestion] = useState<any>();

  const [defaultProgress, setDefaultProgress] = useState(5);
  const [progress, setProgress] = useState(defaultProgress);

  const handleLogin = (quiz, user) => {
    setUsername(user);
    localStorage.setItem('username', user);
    navigate('/' + quiz);
  };

  const handleNewSession = (quiz, user, config) => {
    setUsername(user);
    localStorage.setItem('username', user);
    client.createSesssion(quiz, config);
    setDefaultProgress(config.questionTimer);
    navigate('/' + quiz);
  };

  const handleLogout = (quiz, user) => {
    setUsername(user);
    localStorage.removeItem('username');
    navigate('/');
  };

  const handleChooseAnswer = (questionId, answer) => {
    socket?.send(
      JSON.stringify({
        eventType: 'answer',
        username: username,
        questionId: questionId,
        quizId: quizId,
        answer: answer,
      }),
    );
  };

  useEffect(() => {
    if (quizId && username) {
      const newSocket = new WebSocketReconnector(quizId, username);
      console.log('Create new socket');
      setSocket((prevSocket) => {
        prevSocket?.close();
        return newSocket;
      });
      return () => {
        newSocket.close();
      };
    }
  }, [quizId, username]);

  useEffect(() => {
    socket?.isNotReady() && setOpenReconnecting(true);

    socket?.connect({
      onConnectCallback: () => {
        setOpenReconnecting(false);
      },
      processMessageFn: async (data) => {
        try {
          const json = JSON.parse(data);
          setBoardStatus(json['@']);
          switch (json.eventType) {
            case 'leaderboard':
              setLeaderBoard(json.data);
              break;
            case 'command':
              if (json['@'] == 'starting') {
                let countdownValue = Number(json.data);
                setCounter(countdownValue);
                const interval = setInterval(() => {
                  countdownValue--;
                  setCounter(countdownValue);
                  if (countdownValue == 0) {
                    clearInterval(interval);
                  }
                }, 1000);
              }
              setBoardStatus(json['@']);
              break;
            case 'question':
              setQuestion(json.data);
              break;
          }
        } catch {
          console.log('can not parse JSON message: ', data);
        }
      },
    });
  }, [socket]);

  useEffect(() => {
    setProgress(100);
    const timer = setInterval(() => {
      setProgress((oldProgress) => {
        if (oldProgress < 0) {
          clearInterval(timer);
          return 0;
        }
        const diff = 100 / (defaultProgress - 1);
        return oldProgress - diff;
      });
    }, 1000);
    return () => {
      clearInterval(timer);
    };
  }, [quiz]);

  return (
    <ThemeProvider theme={theme}>
      {username && quizId ? (
        <Container>
          <Stack direction="row" justifyContent="space-between">
            <Typography variant="h5">{username}</Typography>
            <LeaderBoard data={leaderBoard} />
          </Stack>
          <Divider sx={{ marginTop: 2, marginBottom: 2 }} />
          {boardStatus == 'waiting' || boardStatus == 'ended' ? (
            <Button
              variant="contained"
              color="secondary"
              onClick={() => {
                client.startSession(quizId);
              }}
            >
              Start game
            </Button>
          ) : counter > 0 ? (
            <Typography>Starting in ... {counter}</Typography>
          ) : (
            <Typography>Game started</Typography>
          )}
          <Divider sx={{ marginTop: 2, marginBottom: 2 }} />
          {boardStatus == 'inProgress' && quiz?.question && (
            <>
              <LinearProgress
                variant="determinate"
                value={progress}
                color="secondary"
              />
              <QuestionBoard
                key={quiz.id}
                questionId={quiz.id}
                question={quiz.question}
                answers={quiz.answer}
                handleChooseAnswer={handleChooseAnswer}
              />
            </>
          )}
          {boardStatus == 'ended' && (
            <Typography variant="h5" color="error" sx={{ marginTop: 2 }}>
              Game ended
            </Typography>
          )}
        </Container>
      ) : (
        <LoginBox onLogin={handleLogin} onNewSession={handleNewSession} />
      )}
    </ThemeProvider>
  );
}

export default App;
