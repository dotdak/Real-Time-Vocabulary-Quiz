import React from 'react';

import { Container, Paper, Stack, Typography } from '@mui/material';

import { DataGrid, GridColDef } from '@mui/x-data-grid';

interface LeaderBoardProps {
  data: { [key: string]: number };
}

const columns: GridColDef[] = [
  { field: 'id', headerName: 'User', width: 200, sortable: false },
  {
    field: 'points',
    headerName: 'Points',
    width: 200,
    sortable: true,
    sortingOrder: ['desc'],
    renderCell: (params) => {
      return <b>{params.value}</b>;
    },
  },
];

const paginationModel = { page: 0, pageSize: 5 };

export default function LeaderBoard(props: LeaderBoardProps) {
  const rows = Object.entries<number>(props.data)
    .map(([username, points]) => ({
      id: username,
      points,
    }))
    .sort((a, b) => b.points - a.points);
  return (
    <Paper sx={{ width: 400, marginTop: 4 }}>
      <Stack alignItems="center">
        <Typography>Leaderboard</Typography>
        <DataGrid
          rows={rows}
          columns={columns}
          initialState={{ pagination: { paginationModel } }}
          pageSizeOptions={[5, 10]}
          sx={{ border: 0 }}
        />
      </Stack>
    </Paper>
  );
}
