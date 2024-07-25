import { Avatar } from "@nextui-org/avatar";
import { Button } from "@nextui-org/button";
import {
  Dropdown,
  DropdownItem,
  DropdownMenu,
  DropdownTrigger,
} from "@nextui-org/dropdown";
import { Link } from "@nextui-org/link";
import {
  Navbar,
  NavbarBrand,
  NavbarContent,
  NavbarItem,
} from "@nextui-org/navbar";
import { getApp } from "firebase/app";
import {
  GithubAuthProvider,
  getAuth,
  signInWithPopup,
  signOut,
} from "firebase/auth";
import type React from "react";
import { useCallback } from "react";
import { BiLogoGithub as LogoGithub } from "react-icons/bi";
import {
  BiHome as Home,
  BiLogOut as LogOut,
  BiUser as User,
} from "react-icons/bi";
import { navigate } from "vike/client/router";

import logoImg from "@/assets/logo.svg";
import { useFirebaseUser } from "@/hooks/firebase";
import { Image } from "@nextui-org/image";

export default function Layout({
  children,
}: {
  children: React.ReactNode;
}) {
  const fbUser = useFirebaseUser();

  const onSignUpClick = useCallback(async () => {
    await signInWithPopup(getAuth(getApp()), new GithubAuthProvider());
    navigate("/settings");
  }, []);

  const onSettingsClick = useCallback(() => {
    navigate("/settings");
  }, []);

  const onLogOutClick = useCallback(() => {
    if (!fbUser) {
      return;
    }

    signOut(getAuth(getApp()));
  }, [fbUser]);

  return (
    <div className="container mx-auto prose max-w-7xl prose-img:m-0">
      <div className="grid grid-cols-4 gap-6 md:grid-cols-8 lg:grid-cols-12">
        <div className="col-span-4 md:col-span-8 lg:col-span-12">
          <Navbar>
            <NavbarBrand>
              <Link href="/" className="text-content0">
                <Image
                  classNames={{ wrapper: "inline-block align-middle mr-2" }}
                  src={logoImg}
                />
                <p className="inline-block align-middle font-extrabold">
                  TASUKE
                </p>
              </Link>
            </NavbarBrand>
            {fbUser ? (
              <>
                <NavbarContent as="div" justify="end">
                  <Dropdown placement="bottom-end">
                    <DropdownTrigger>
                      <Avatar src={fbUser.photoURL ?? undefined} />
                    </DropdownTrigger>
                    <DropdownMenu aria-label="User Actions" variant="flat">
                      <DropdownItem
                        key="settings"
                        onClick={onSettingsClick}
                        startContent={<User className="mr-2 h-4 w-4" />}
                      >
                        Settings
                      </DropdownItem>
                      <DropdownItem
                        key="logout"
                        onClick={onLogOutClick}
                        startContent={<LogOut className="mr-2 h-4 w-4" />}
                      >
                        Log out
                      </DropdownItem>
                    </DropdownMenu>
                  </Dropdown>
                </NavbarContent>
              </>
            ) : (
              <NavbarContent justify="end">
                <NavbarItem>
                  <Button
                    className="bg-primary-400 text-content1"
                    onClick={onSignUpClick}
                    startContent={<LogoGithub className="size-6" />}
                  >
                    Sign Up
                  </Button>
                </NavbarItem>
              </NavbarContent>
            )}
          </Navbar>
        </div>
        {children}
      </div>
    </div>
  );
}
