/* eslint-disable react-refresh/only-export-components */
import React from "react";
import { useLocation, useNavigate } from "react-router-dom";

type UserData = {
  address: string;
  chainId: number;
  exp?: number;
};

type AuthContextType = {
  user: UserData | null;
  jwt: string | null;
  signIn: (token: string) => void;
  signOut: () => void;
};

export const AuthContext = React.createContext<AuthContextType>({
  user: null,
  jwt: null,
  signIn: () => {},
  signOut: () => {},
} as AuthContextType);

const extractUserData = (token: string) => {
  if (!token) return null;

  const [, body] = token.split(".");
  const payload = JSON.parse(atob(body));

  return {
    exp: payload.exp,
    address: payload.sub,
    chainId: Number(payload.cId),
  };
};

export const AuthProvider: React.FC<{ children: React.ReactNode }> = ({
  children,
}) => {
  const navigate = useNavigate();
  const location = useLocation();

  const [user, setUser] = React.useState<UserData | null>(
    () => extractUserData(localStorage.getItem("jwt") || "") || null
  );
  const [jwt, setJwt] = React.useState<string | null>(() =>
    localStorage.getItem("jwt")
  );
  const expirationTimeout = React.useRef<number | null>(null);

  const clearExpirationTimeout = () => {
    if (expirationTimeout.current != null) {
      window.clearTimeout(expirationTimeout.current);
      expirationTimeout.current = null;
    }
  };

  const signIn = React.useCallback((token: string) => {
    localStorage.setItem("jwt", token);
    setJwt(token);
  }, []);

  const signOut = React.useCallback(() => {
    clearExpirationTimeout();
    localStorage.removeItem("jwt");
    setJwt(null);
    setUser(null);
  }, []);

  React.useEffect(() => {
    clearExpirationTimeout();

    if (!jwt) {
      setUser(null);
      return undefined;
    }

    const decoded = extractUserData(jwt);
    setUser(decoded);

    if (!decoded?.exp) {
      return undefined;
    }

    const now = Date.now();
    const expiresAtMs = decoded.exp * 1000;

    if (expiresAtMs <= now) {
      signOut();
      return undefined;
    }

    const delayMs = expiresAtMs - now;
    expirationTimeout.current = window.setTimeout(signOut, delayMs);

    return () => clearExpirationTimeout();
  }, [jwt, signOut]);

  React.useEffect(() => {
    const isLoginRoute = location.pathname === "/login";
    const isVerifyWalletRoute = location.pathname === "/verify-wallet";

    if (jwt && isLoginRoute) {
      const params = new URLSearchParams(window.location.search);
      const redirect = params.get("redirect");

      if (redirect) {
        console.log("Redirecting to:", redirect);
        navigate(redirect);
        return;
      }

      navigate("/dashboard");
    } else if (!jwt && !isLoginRoute) {
      const params = new URLSearchParams();
      if (isVerifyWalletRoute) {
        params.set("redirect", location.pathname + location.search);
      }

      navigate(`/login?${params.toString()}`);
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [jwt, location]);

  return (
    <AuthContext.Provider value={{ user, jwt, signIn, signOut }}>
      {children}
    </AuthContext.Provider>
  );
};

export const useAuth = () => {
  const context = React.useContext(AuthContext);
  if (!context) throw new Error("useAuth must be used within an AuthProvider");
  return context;
};
