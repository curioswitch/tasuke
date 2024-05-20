import { initializeApp } from "firebase/app";
import { type User, getAuth, onAuthStateChanged } from "firebase/auth";
import { createContext, useContext, useEffect, useState } from "react";

import { getFirebaseConfig } from "./config";

const app = initializeApp(getFirebaseConfig());
const auth = getAuth(app);

interface FirebaseState {
  user?: User;
  userResolved?: boolean;
}

const FirebaseContext = createContext<FirebaseState | undefined>(undefined);

export function useFirebaseState(): FirebaseState | undefined {
  return useContext(FirebaseContext);
}

export function useFirebaseUser(): User | undefined {
  return useContext(FirebaseContext)?.user;
}

export function FirebaseProvider({ children }: { children: React.ReactNode }) {
  const [state, setState] = useState<FirebaseState | undefined>(undefined);

  useEffect(() => {
    return onAuthStateChanged(auth, (u) => {
      const user = u ?? undefined;
      setState({ user, userResolved: true });
    });
  }, []);

  return (
    <FirebaseContext.Provider value={state}>
      {children}
    </FirebaseContext.Provider>
  );
}
