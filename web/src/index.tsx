import React from 'react';
import ReactDOM from 'react-dom/client';
import './index.css';
import App from './App';
import { BrowserRouter, Route, Routes } from 'react-router-dom';

const root = ReactDOM.createRoot(
  document.getElementById('root') as HTMLElement,
);

root.render(
  <React.StrictMode>
    <div>
      <BrowserRouter>
        <Routes>
          <Route path="/" element={<App />}>
            <Route path=":quizId" element={<App />} />
          </Route>
        </Routes>
      </BrowserRouter>
    </div>
  </React.StrictMode>,
);
