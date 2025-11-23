/* eslint-disable react-refresh/only-export-components */
import React from "react";
import { Eip1193Provider } from "ethers";
import { useAppKitProvider } from "@reown/appkit/react";

type WalletProviderContextType = {
  wallet: Eip1193Provider | null;
};

export const WalletProviderContext =
  React.createContext<WalletProviderContextType>(
    {} as WalletProviderContextType
  );

export const WalletProvider: React.FC<{ children: React.ReactNode }> = ({
  children,
}) => {
  const { walletProvider } = useAppKitProvider("eip155");
  const eip1193Provider = walletProvider as Eip1193Provider;

  return (
    <WalletProviderContext.Provider value={{ wallet: eip1193Provider }}>
      {children}
    </WalletProviderContext.Provider>
  );
};

export const useWallet = () => {
  const context = React.useContext(WalletProviderContext);
  if (!context)
    throw new Error("useWallet must be used within a WalletProvider");
  return context;
};
