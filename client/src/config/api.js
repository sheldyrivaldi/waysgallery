import axios from "axios";

// Create Base URL
// const baseURL = import.meta.env.VITE_BASE_URL;
const baseURL = "http://localhost:5000/api/v1";
export const API = axios.create({
  baseURL: baseURL,
});

// Authorization Token
export const setAuthToken = (token) => {
  if (token) {
    API.defaults.headers.common["Authorization"] = `Bearer ${token}`;
  } else {
    delete API.defaults.headers.common["Authorization"];
  }
};
