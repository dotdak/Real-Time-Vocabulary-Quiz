import React, { useState } from 'react';
import styles from './LoginBox.module.css'; // Import the CSS file for styling
import { useParams } from 'react-router-dom';
import {
  Button,
  Divider,
  InputAdornment,
  Stack,
  Grid,
  TextField,
  Typography,
  Box,
} from '@mui/material';

import AccountCircle from '@mui/icons-material/AccountCircle';
import { MeetingRoomRounded } from '@mui/icons-material';

interface LoginBoxProps {
  onLogin: (quizId, username) => void;
  onNewSession: (quizId, username, config) => void;
}
const LoginBox = ({ onLogin, onNewSession }: LoginBoxProps) => {
  const { quizId: quizId } = useParams();

  const [isLogin, setIsLogin] = useState(true);
  const [username, setUsername] = useState('');
  const [_quizId, setQuizId] = useState(quizId || '');
  const [totalQuestion, setTotalQuestion] = useState(10);
  const [questionTimer, setQuestionTimer] = useState(5);

  const handleSessionInputChange = (event) => {
    setQuizId(event.target.value);
  };

  const handleUserInputChange = (event) => {
    setUsername(event.target.value);
  };

  const handleKeyUp = (event) => {
    if (event.key === 'Enter') {
      handleLogin();
    }
  };

  const validateInput = (input: string): boolean => {
    return (
      input !== '@' &&
      input.trim() !== '' &&
      !input.includes(':') &&
      !input.includes('/')
    );
  };

  const handleLogin = () => {
    if (validateInput(username) && validateInput(_quizId))
      onLogin(_quizId, username);
  };

  const handleNewSession = () => {
    if (validateInput(username) && validateInput(_quizId))
      onNewSession(_quizId, username, {
        questionTimer: questionTimer * 1000000000,
        totalQuestion,
      });
  };

  let content = (
    <Stack alignItems="center">
      <TextField
        id="input-room"
        label="Session ID"
        value={_quizId}
        onChange={handleSessionInputChange}
        onKeyUp={(event) => {
          if (event.key === 'Enter') {
            handleNewSession();
          }
        }}
        InputProps={{
          startAdornment: (
            <InputAdornment position="start">
              <MeetingRoomRounded />
            </InputAdornment>
          ),
        }}
        variant="standard"
        sx={{ marginY: '1rem', width: 200 }}
      />
      <TextField
        id="input-username"
        label="Username"
        value={username}
        onChange={handleUserInputChange}
        onKeyUp={(event) => {
          if (event.key === 'Enter') {
            handleNewSession();
          }
        }}
        InputProps={{
          startAdornment: (
            <InputAdornment position="start">
              <AccountCircle />
            </InputAdornment>
          ),
        }}
        variant="standard"
        sx={{ marginY: '1rem', width: 200 }}
      />
      <Divider sx={{ marginBottom: '1rem' }} />
      <Grid container width={400} spacing={2} columns={4}>
        <Grid item xs={2} justifyItems="start">
          <Typography>Number of questions:</Typography>
        </Grid>
        <Grid item xs={1}>
          <TextField
            id="outlined-number"
            type="number"
            size="small"
            value={totalQuestion}
            onChange={(event) => {
              let value = Number(event.target.value);
              if (value < 0) {
                value = 0;
              }

              setTotalQuestion(value);
            }}
          />
        </Grid>
        <Grid item xs={1}>
          <Typography>questions</Typography>
        </Grid>
        <Grid item xs={2} justifyItems="start">
          <Typography>Timer:</Typography>
        </Grid>
        <Grid item xs={1}>
          <TextField
            id="outlined-number"
            type="number"
            size="small"
            value={questionTimer}
            onChange={(event) => {
              let value = Number(event.target.value);
              if (value < 0) {
                value = 0;
              }
              setQuestionTimer(value);
            }}
          />
        </Grid>
        <Grid item xs={1}>
          <Typography>seconds</Typography>
        </Grid>
      </Grid>
      <br />
      <Button
        onClick={handleNewSession}
        sx={{ marginY: '1rem' }}
        variant="contained"
      >
        Create session
      </Button>
    </Stack>
  );
  if (isLogin)
    content = (
      <>
        {quizId ? (
          <Typography variant="body1">
            Join room{' '}
            <Typography variant="overline" color="secondary">
              {quizId}
            </Typography>{' '}
            ?
          </Typography>
        ) : (
          <>
            <TextField
              id="input-room"
              label="Session ID"
              value={_quizId}
              onChange={handleSessionInputChange}
              onKeyUp={handleKeyUp}
              InputProps={{
                startAdornment: (
                  <InputAdornment position="start">
                    <MeetingRoomRounded />
                  </InputAdornment>
                ),
              }}
              variant="standard"
              sx={{ marginY: '1rem' }}
            />
            <br />
          </>
        )}
        <TextField
          id="input-username"
          label="Username"
          value={username}
          onChange={handleUserInputChange}
          onKeyUp={handleKeyUp}
          InputProps={{
            startAdornment: (
              <InputAdornment position="start">
                <AccountCircle />
              </InputAdornment>
            ),
          }}
          variant="standard"
          sx={{ marginY: '1rem' }}
        />
        <br />
        <Button
          onClick={handleLogin}
          sx={{ marginY: '1rem' }}
          variant="contained"
        >
          Login
        </Button>
        <br />
        <Button
          onClick={() => {
            setIsLogin(false);
          }}
          sx={{ marginY: '1rem' }}
          variant="outlined"
          color="secondary"
        >
          Start a new session ?
        </Button>
      </>
    );

  return (
    <div className={styles.loginContainer}>
      <Stack>
        <Typography variant="h2" color="primary.dark">
          Vocabulary Quiz
        </Typography>
        <div className={styles.loginBox}>{content}</div>
      </Stack>
    </div>
  );
};

export default LoginBox;
