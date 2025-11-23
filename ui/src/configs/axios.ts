import axios from "axios";

export const initAxiosDefaults = () => {
  axios.defaults.baseURL =
    import.meta.env.VITE_ORCHESTRATOR_GATEWAY_API_BASE_URL;
  axios.defaults.headers.post["Content-Type"] = "application/json";
};
