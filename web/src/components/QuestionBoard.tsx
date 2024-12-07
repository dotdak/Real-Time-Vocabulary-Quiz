import React, { useState } from 'react';
import {
  Button,
  Card,
  CardActions,
  CardContent,
  Typography,
} from '@mui/material';

export interface QuestionBoardProps {
  questionId: string;
  question: string;
  answers: { [key: string]: string };
  handleChooseAnswer: (questionId: string, answerKey: string) => void;
}

export default function QuestionBoard(props: QuestionBoardProps) {
  const [selected, setSelected] = useState<string>();
  return (
    <Card>
      <CardContent>
        <Typography fontWeight="600">{props.question}</Typography>
      </CardContent>
      <CardActions>
        {Object.entries(props.answers).map(([answerKey, answerValue]) => (
          <Button
            key={answerKey}
            size="small"
            color="primary"
            variant="contained"
            disabled={selected != undefined}
            onClick={() => {
              setSelected(answerKey);
              props.handleChooseAnswer(props.questionId, answerKey);
            }}
          >
            {answerKey}. {answerValue}
          </Button>
        ))}
      </CardActions>
    </Card>
  );
}
