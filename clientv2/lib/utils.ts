import { type ClassValue, clsx } from 'clsx'
import { twMerge } from 'tailwind-merge'

/**
 * Combines multiple class values into a single className string.
 * Uses clsx for conditional classes and tailwind-merge to handle Tailwind CSS class conflicts.
 * 
 * @param inputs - Class values to be combined (strings, objects, arrays, etc.)
 * @returns A merged className string
 */
export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs))
}