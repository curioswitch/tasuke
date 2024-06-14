import { Button } from "@nextui-org/button";
import { getApp } from "firebase/app";
import { GithubAuthProvider, getAuth, signInWithPopup } from "firebase/auth";
import { useCallback } from "react";
import { BiLogoGithub as LogoGithub } from "react-icons/bi";
import { navigate } from "vike/client/router";

import { useFirebase } from "@/hooks/firebase";

function getBotInstallLink() {
  if (import.meta.env.PUBLIC_ENV__FIREBASE_APP === "tasuke-dev") {
    return "https://github.com/apps/tasuke-alpha";
  }
  if (import.meta.env.PUBLIC_ENV__FIREBASE_APP === "tasuke-prod") {
    return "https://github.com/apps/tasuke-bot";
  }
  throw new Error("PUBLIC_ENV__FIREBASE_APP must be configured");
}

export default function Page() {
  const fbState = useFirebase();

  const onSignUpClick = useCallback(async () => {
    await signInWithPopup(getAuth(getApp()), new GithubAuthProvider());
    navigate("/settings");
  }, []);

  return (
    <>
      <div className="col-span-4 md:col-span-8 lg:col-span-12">
        <h1>Tasuke</h1>
        <p className="lead">
          tasuke is short for the Japanese 手助け (tedasuke), meaning to give a
          helping hand. It aims to support OSS developers in giving each other a
          helping hand by connecting code reviewers to PRs in otherwise
          unrelated OSS repositories.
        </p>
        <p>
          Everyone knows the{" "}
          <a
            href="https://xkcd.com/2347/"
            className="underline"
            target="_blank"
            rel="noreferrer noopener"
          >
            lone maintainer
          </a>{" "}
          issue in OSS - let's see if we can make it at least a little less
          lonely.
        </p>
        <h2>Getting Started</h2>
        <h3>For code reviewers</h3>
        <p>
          {fbState?.userResolved && !fbState.user ? (
            <>
              <p>
                If you are interested in helping with code reviews, create an
                account and submit your preferences for reviewing.
              </p>
              <Button
                className="bg-primary-500 text-content1"
                onClick={onSignUpClick}
                startContent={<LogoGithub className="size-6" />}
              >
                Sign up with GitHub
              </Button>
            </>
          ) : (
            <>
              <p>
                If you are interested in helping with code reviews, submit or
                edit your preferences for reviewing.
              </p>
              <p>
                <a href="/settings">Edit your preferences</a>
              </p>
            </>
          )}
        </p>
        <h3>For maintainers</h3>
        <p>
          If you have a repository you could use help with on code review,
          register the Tasuke GitHub app.
        </p>
        <blockquote>
          Currently you must also register as a code reviewer to use the app. As
          review capacity increases and we understand abuse potential better, we
          may remove this requirement.
        </blockquote>
        <p>
          <a
            href={getBotInstallLink()}
            target="_blank"
            rel="noreferrer noopener"
          >
            Install the app
          </a>
        </p>
        <p>
          With the app installed on a repository, create a pull request in it
          and request a review by adding a comment starting with{" "}
          <code>/tasuke</code>, for example:
        </p>
        <p>
          <code>/tasuke Can you help me with a review, please?</code>
        </p>
        <p>
          Currently, the only supported command is to ask for a review, so the
          text content is not used. In the future, we may use it for executing
          different commands.
        </p>
        <h3>Contributing to tasuke</h3>
        <p>tasuke is OSS for OSS.</p>
        <p>
          <a
            href="https://github.com/curioswitch/tasuke"
            target="_blank"
            rel="noreferrer noopener"
          >
            Check out the repo
          </a>
        </p>
        <p>
          If you are interested, feel free to file issues or pull requests as
          with any OSS project. Let's work together to find what can make the
          service most useful for OSS developers.
        </p>
      </div>
    </>
  );
}
