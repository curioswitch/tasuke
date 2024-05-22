import { Avatar } from "@nextui-org/avatar";
import {
  Dropdown,
  DropdownItem,
  DropdownMenu,
  DropdownTrigger,
} from "@nextui-org/dropdown";
import { getApp } from "firebase/app";
import { getAuth, signOut } from "firebase/auth";
import type React from "react";
import { useCallback } from "react";
import { SlLogout, SlUser } from "react-icons/sl";
import { navigate } from "vike/client/router";

import { useFirebaseUser } from "@/hooks/firebase";

export default function Layout({
  children,
}: {
  children: React.ReactNode;
}) {
  const fbUser = useFirebaseUser();

  // href not working for some reason.
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
    <>
      <section>
        <div className="container px-6 pt-6 mx-auto">
          <div className="grid grid-cols-4 gap-6 md:grid-cols-8 lg:grid-cols-12 prose lg:prose-xl">
            <div className="flex justify-end col-span-4 md:col-span-8 lg:col-span-12">
              {fbUser ? (
                <Dropdown>
                  <DropdownTrigger className="cursor-pointer">
                    <Avatar src={fbUser.photoURL ?? undefined} />
                  </DropdownTrigger>
                  <DropdownMenu aria-label="User Actions">
                    <DropdownItem
                      key="settings"
                      onClick={onSettingsClick}
                      startContent={<SlUser className="mr-2 h-4 w-4" />}
                    >
                      Settings
                    </DropdownItem>
                    <DropdownItem
                      key="logout"
                      onClick={onLogOutClick}
                      startContent={<SlLogout className="mr-2 h-4 w-4" />}
                    >
                      Log out
                    </DropdownItem>
                  </DropdownMenu>
                </Dropdown>
              ) : (
                <Avatar src={undefined} />
              )}
            </div>
            {children}
          </div>
        </div>
      </section>
    </>
  );
}
