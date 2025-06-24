# MoodMap — Real-Time Geospatial Mood Dashboard

**MoodMap** is a full-stack web application that visualizes user-submitted mood emojis across India in real-time. It uses **Go**, **PostgreSQL**, and **Leaflet.js** to cluster, anonymize, and display mood data as heatmap blobs based on geospatial patterns.

## Project Overview

MoodMap allows users to post their current mood as an emoji along with an optional text. These moods are plotted on a map with real-time updates and privacy-aware location rounding. Using DBSCAN clustering and convex hull visualization, the app reveals emerging emotional hotspots and trends across the country.

## Tech Stack

- **Backend:** Go (Golang), REST APIs  
- **Database:** PostgreSQL with PostGIS (for geospatial indexing)  
- **Frontend:** HTML/CSS + Leaflet.js for interactive maps  
- **Clustering & Visualization:** DBSCAN (density-based) + convex hulls + color-coded heatmap blobs

## Key Features

- **Emoji Mood Posts:** Users post their mood using emojis and optional messages  
- **Location-Aware Visualization:** Approximate geolocation (rounded to ~100km) preserves privacy while enabling clustering  
- **Dynamic Heatmap Blobs:** DBSCAN clusters moods based on proximity (≤ 500km), displayed using colored convex hulls  
- **Real-Time Updates:** New moods appear instantly on the map  
- **Privacy Protection:** Coordinates are rounded to obscure exact user location while retaining usefulness for regional mood trends

## Clustering & Visualization

- **Algorithm:** DBSCAN clusters mood points using haversine distance  
- **Rendering:** Clusters are rendered as convex hulls on the Leaflet map, with:
  - Color based on mood emoji  
  - Hover tooltips showing emoji counts and top messages
