import vikeReact from "vike-react/config";
import type { Config } from "vike/types";

import Head from "@/layouts/Head.jsx";
import Layout from "@/layouts/Layout.jsx";
import Wrapper from "@/layouts/Wrapper.jsx";

// Default config (can be overridden by pages)
export default {
  Layout,
  Head,
  Wrapper,
  title: "tasuke",
  extends: vikeReact,
} satisfies Config;
