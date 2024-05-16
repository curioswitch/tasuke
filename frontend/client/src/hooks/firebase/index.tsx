import { type FirebaseApp, initializeApp } from "firebase/app";
import { type Auth, type User, getAuth } from "firebase/auth";
import { createContext, useContext, useEffect, useState } from "react";

interface FirebaseState {
  app: FirebaseApp;
  auth: Auth;
  user?: User;
}

const FirebaseContext = createContext<FirebaseState | undefined>(undefined);

export function useAuth(): Auth | undefined {
  return useContext(FirebaseContext)?.auth;
}

export function useUser(): User | undefined {
  return useContext(FirebaseContext)?.user;
}

export function FirebaseProvider({ children }: { children: React.ReactNode }) {
  const [state, setState] = useState<FirebaseState | undefined>(undefined);

  useEffect(() => {
    async function fetchConfig() {
      const response = await fetch("/__/firebase/init.json");
      const config = await response.json();
      const app = initializeApp(config);
      const auth = getAuth(app);
      setState({ app, auth });

      auth.onAuthStateChanged((u) => {
        const user = u ?? undefined;
        setState({ app, auth, user });
      });
    }

    fetchConfig();
  }, []);

  return (
    <FirebaseContext.Provider value={state}>
      {children}
    </FirebaseContext.Provider>
  );
}
