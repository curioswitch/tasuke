export function getFirebaseConfig() {
  // Note, switch statement does not seem to allow optimizing away the unused
  // config so we use if statements.

  if (import.meta.env.PUBLIC_ENV__FIREBASE_APP === "tasuke-dev") {
    return {
      apiKey: "AIzaSyDjdtrbYI3kYWd0YUHMSPLeXevjHyZGGlY",
      authDomain: "alpha.tasuke.dev",
      projectId: "tasuke-dev",
      storageBucket: "tasuke-dev.appspot.com",
      messagingSenderId: "720364425367",
      appId: "1:720364425367:web:509f4c126ae54228bfd9d2",
    };
  }

  if (import.meta.env.PUBLIC_ENV__FIREBASE_APP === "tasuke-prod") {
    return {
      apiKey: "AIzaSyDtWCei1awTdnfWIE9yrCU9PhqOy5qqJ9w",
      authDomain: "tasuke.dev",
      projectId: "tasuke-prod",
      storageBucket: "tasuke-prod.appspot.com",
      messagingSenderId: "840011577241",
      appId: "1:840011577241:web:db8a0a45a044be541dc508",
    };
  }

  throw new Error("PUBLIC_ENV__FIREBASE_APP must be configured");
}
