# Social Media Admin Dashboard

A modern, high-performance admin dashboard for social media platform management. Built with **Next.js 15**, **Tailwind CSS v4**, and a custom **Glassmorphism** design system.

## Key Features

- **Monochromatic Glassmorphism UI**
  - Unified design language using `backdrop-blur-xl` and semi-transparent backgrounds.
  - Consistent use of `bg-background/60` and `bg-card/50` for a depth-rich interface.
  - Fully responsive layout with mobile-optimized navigation.

- **Command Center Dashboard**
  - Real-time statistics cards with hover effects.
  - Interactive charts using Recharts.
  - Live activity feed simulation.
  - "Quick Actions" and "Moderation Queue" shortcuts.

- **Comprehensive Settings**
  - Tabbed interface for General, Notifications, and Platform controls.
  - Functional UI for profile management and system toggles.
  - Responsive design that adapts to mobile screens.

- **Interactive Experience**
  - Toast notifications via `sonner` for user feedback.
  - Micro-interactions (hover scales, transitions) on all interactive elements.
  - Client-side routing with immediate redirects (Root `/` redirects to `/dashboard`).

## Tech Stack

- **Framework**: [Next.js](https://nextjs.org/) (App Router)
- **Styling**: [Tailwind CSS](https://tailwindcss.com/)
- **UI Components**: [Shadcn UI](https://ui.shadcn.com/) (Radix Primitives)
- **Icons**: [Lucide React](https://lucide.dev/)
- **Charts**: [Recharts](https://recharts.org/)
- **State Management**: [TanStack Query](https://tanstack.com/query/latest)
- **Notifications**: [Sonner](https://sonner.emilkowal.ski/)

## Getting Started

1.  **Install dependencies**:
    ```bash
    npm install
    ```

2.  **Run the development server**:
    ```bash
    npm run dev
    ```

3.  **Open the application**:
    Navigate to [http://localhost:3000](http://localhost:3000). The application will automatically redirect you to the dashboard.

## Project Structure

```
frontend/
├── app/
│   ├── dashboard/       # Dashboard routes and pages
│   │   ├── settings/    # Settings page
│   │   └── page.tsx     # Main Command Center
│   ├── globals.css      # Global styles & Tailwind theme
│   ├── layout.tsx       # Root layout with providers
│   └── page.tsx         # Root redirect
├── components/
│   ├── dashboard/       # Dashboard-specific components
│   ├── ui/              # Reusable UI components (Shadcn)
│   ├── header.tsx       # Top navigation bar
│   └── sidebar.tsx      # Side navigation
├── lib/                 # Utilities and API mocks
└── public/              # Static assets
```

## Design System

The project uses a custom Tailwind configuration to achieve the glassmorphism look:

- **Backgrounds**: Heavy use of alpha channels (e.g., `bg-background/60`).
- **Blur**: `backdrop-blur-xl` is applied globally to cards, sidebars, and headers.
- **Borders**: Subtle borders (`border-border/40`) to define edges without heavy lines.
- **Shadows**: Soft shadows that increase on hover for depth perception.
