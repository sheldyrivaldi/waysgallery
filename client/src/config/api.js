import axios from "axios";

// Create Base URL
const baseURL = import.meta.env.VITE_BASE_URL;
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
