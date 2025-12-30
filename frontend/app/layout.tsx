import { Inter } from "next/font/google";

import type { Metadata } from "next";
import "./globals.css";

const inter = Inter({ subsets: ["latin"] });

export const metadata: Metadata = {
  title: "CI/CD Training Application",
  description: "Learn DevOps practices through hands-on CI/CD exercises",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en">
      <body className={inter.className}>
        <div className="min-h-screen bg-background">
          <header className="border-b">
            <div className="container mx-auto px-4 py-4">
              <nav className="flex items-center justify-between">
                <div className="flex items-center space-x-8">
                  <a href="/" className="text-xl font-bold">
                    CI/CD Training
                  </a>
                  <div className="flex space-x-4">
                    <a
                      href="/"
                      className="text-sm font-medium text-muted-foreground hover:text-foreground transition-colors"
                    >
                      Dashboard
                    </a>
                    <a
                      href="/health"
                      className="text-sm font-medium text-muted-foreground hover:text-foreground transition-colors"
                    >
                      Health
                    </a>
                    <a
                      href="/deployments"
                      className="text-sm font-medium text-muted-foreground hover:text-foreground transition-colors"
                    >
                      Deployments
                    </a>
                  </div>
                </div>
                <div className="flex items-center space-x-4">
                  <span className="text-sm text-muted-foreground">
                    Environment: {process.env.NODE_ENV || "development"}
                  </span>
                </div>
              </nav>
            </div>
          </header>
          <main className="container mx-auto px-4 py-8">
            {children}
          </main>
          <footer className="border-t mt-auto">
            <div className="container mx-auto px-4 py-6">
              <p className="text-sm text-center text-muted-foreground">
                CI/CD Training Application - Learn DevOps Practices
              </p>
            </div>
          </footer>
        </div>
      </body>
    </html>
  );
}
