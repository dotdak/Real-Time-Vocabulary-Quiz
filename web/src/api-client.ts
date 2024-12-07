import axios, { Axios } from 'axios';

const API_URL = process.env.REACT_APP_API_URI || '';

export class ApiClient {
  client: Axios;

  constructor() {
    this.client = axios.create({
      baseURL: API_URL,
    });
    this.client.interceptors.response.use(
      (response) => response,
      (error) => console.error(error),
    );
  }

  async createSesssion(quizId, config) {
    await this.client.post('/api/quiz', {
      quizId,
      config,
    });
  }

  async startSession(quizId: string) {
    await this.client.put(`/api/quiz/${quizId}`, {
      status: 'inProgress',
    });
  }
}

const client = new ApiClient();

export default client;
