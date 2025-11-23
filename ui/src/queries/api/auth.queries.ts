import axios from "axios";

export const initJwtAuth = async (address: string) => {
  return await axios.post(
    `/auth/jwt/init`,
    { address },
    { headers: { "Content-Type": "application/json" } }
  );
};

export const finishJwtAuth = async (message: string, signature: string) => {
  return await axios.post(
    `/auth/jwt/finish`,
    { message, signature },
    { headers: { "Content-Type": "application/json" } }
  );
};
