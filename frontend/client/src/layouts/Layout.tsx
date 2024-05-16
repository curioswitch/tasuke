import type React from "react";

export default function Layout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <section>
      <div className="container px-6 m-auto">
        <div className="grid grid-cols-4 gap-6 md:grid-cols-8 lg:grid-cols-12">
          {children}
        </div>
      </div>
    </section>
  );
}
