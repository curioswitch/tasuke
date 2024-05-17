import { useQuery } from "@tanstack/react-query";
import { type FirebaseApp, initializeApp } from "firebase/app";
import { type Auth, type User, getAuth } from "firebase/auth";
import { createContext, useContext, useEffect, useState } from "react";

interface FirebaseState {
  app: FirebaseApp;
  auth: Auth;
  user?: User;
  userResolved?: boolean;
}

const FirebaseContext = createContext<FirebaseState | undefined>(undefined);

export function useFirebase(): FirebaseState | undefined {
  return useContext(FirebaseContext);
}

export function useFirebaseUser(): User | undefined {
  return useContext(FirebaseContext)?.user;
}

async function fetchFirebaseConfig() {
  const response = await fetch("/__/firebase/init.json");
  if (!response.ok) {
    throw new Error("Failed to fetch Firebase config");
  }
  const config = await response.json();

  if (import.meta.env.MODE === "development") {
    config.authDomain = "alpha.tasuke.dev";
  } else {
    config.authDomain = window.location.hostname;
  }
  return config;
}

export function FirebaseProvider({ children }: { children: React.ReactNode }) {
  const [state, setState] = useState<FirebaseState | undefined>(undefined);

  const { data, error } = useQuery({
    queryKey: ["firebaseConfig"],
    queryFn: fetchFirebaseConfig,
    // Firebase config is static
    refetchOnMount: false,
    refetchOnWindowFocus: false,
  });

  useEffect(() => {
    if (error) {
      // Extremely unlikely
      // TODO: Handle error
      return;
    }
    if (!data) {
      // Still loading
      return;
    }

    const app = initializeApp(data);
    const auth = getAuth(app);
    setState({ app, auth });

    auth.onAuthStateChanged((u) => {
      const user = u ?? undefined;
      setState({ app, auth, user, userResolved: true });
    });
  }, [data, error]);

  return (
    <FirebaseContext.Provider value={state}>
      {children}
    </FirebaseContext.Provider>
  );
}
