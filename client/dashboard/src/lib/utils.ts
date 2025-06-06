import { clsx, type ClassValue } from "clsx";
import { twMerge } from "tailwind-merge";

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs));
}

// TODO: Will potentitally map to a specific separate API URL if we deploy through Vercel
export function getServerURL(): string {
  if (__GRAM_SERVER_URL__) {
    return __GRAM_SERVER_URL__;
  }

  return window.location.origin;
}

export function titleCase(str: string) {
  return str.replace(/\b\w/g, (char) => char.toUpperCase());
}

export function capitalize(str: string) {
  return str.charAt(0).toUpperCase() + str.slice(1);
}

export function assert(condition: unknown, message: string): asserts condition {
  if (!condition) {
    throw new Error(message);
  }
}
