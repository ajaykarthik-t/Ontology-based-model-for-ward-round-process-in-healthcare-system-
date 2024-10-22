# Ontology-Based Model for Ward Round Process in Healthcare System

This repository contains an ontology-based model designed to improve the ward round process in healthcare systems. The model utilizes an efficient framework to track whether doctors have checked patients or not.

## Features

- **Doctor Management**: CRUD operations for managing doctor information.
- **Patient Management**: CRUD operations for managing patient details.
- **Check Status**: Track whether a patient has been checked by a doctor.

## Technologies Used

- Go
- MongoDB
- Fiber (Web Framework for Go)

## API Endpoints

### Base URL

#### 1. Get All Doctors

- **Method**: `GET`
- **Endpoint**: `/doctor`
- **Description**: Retrieve a list of all doctors.

**Example Request**:
```http
GET /doctor HTTP/1.1
Host: localhost:3000

