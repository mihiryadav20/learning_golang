# Project Documentation: Sports Tryouts API

## Overview
We've been building a backend API service using Go and Fiber for a sports tryouts platform called "Tryouts". This platform allows sports teams to register, create tryout events, and lets users search for and view upcoming tryouts.

## Technical Stack
- **Language**: Go
- **Web Framework**: Fiber
- **Database**: SQLite with GORM (Object-Relational Mapper)
- **Authentication**: JWT (JSON Web Tokens)

## Core Features Implemented

### Team Management
- Team registration with email/password
- Secure authentication using JWT tokens
- Login/logout functionality
- Team profile management

### Tryout Management
- Creating new tryouts with details like title, description, dates, location
- Updating existing tryouts
- Deleting tryouts
- Viewing all tryouts created by a team

### Public Tryout Access
- Searching tryouts by league, division, and location
- Viewing upcoming tryouts (limited to 10, sorted by date)
- Viewing a specific tryout by ID

## API Endpoints

### Authentication
- `POST /api/teams/register` - Register a new team
- `POST /api/teams/login` - Log in and receive a JWT token
- `POST /api/auth/logout` - Log out (requires authentication)
- `GET /api/auth/me` - Get current authenticated team info

### Tryout Management (Protected)
- `POST /api/auth/tryouts` - Create a new tryout
- `GET /api/auth/tryouts/my` - Get all tryouts for the current team
- `PUT /api/auth/tryouts/:id` - Update a tryout
- `DELETE /api/auth/tryouts/:id` - Delete a tryout

### Public Tryout Access
- `GET /api/tryouts/upcoming` - Get upcoming tryouts
- `GET /api/tryouts/search` - Search tryouts with filters
- `GET /api/tryouts/:id` - Get a specific tryout by ID

## Data Models

### Team
Represents a sports team that can register on the platform:
- ID, email, password (hashed), team name, description, etc.

### Tryout
Represents a sports tryout event:
- ID, team ID, title, description, location
- Dates: tryout date, start/end dates, last registration date
- League, division, form link

## Current Status
We have implemented all the core functionality of the API, including:
- User authentication with JWT
- CRUD operations for tryouts
- Search functionality
- Data validation

The application is ready for basic testing and can be extended with additional features in the future.
