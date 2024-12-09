# Cars viewer

## How to run

Navigate to the api directory and input the command make run to start up the api
```bash
cd api
make run
```

After that, open up a new terminal and type the command
```bash
go run .
```
to start the server.  

Finally, navigate on your browser to 
```
http://localhost:8080
```
in order to view the page.

## Page overview

### Home

Here you can apply filters to show only the car models which align with the selected filters.

### Models

A list of all car models in the database. Click on a model to navigate to it's details page.

### Details

Detailed information about a car model.

### Manufacturers

Info about the manufacturers present in the database.

### Comparison

Select two car models for comparison.

# Cars API

This is a simple cars API backed with NodeJS and Express. The app runs on localhost and uses port 3000.
At any time you can change the port to whatever you want, see [Instructions](##Instructions).

## Installation of required packages
- Install [NodeJS](http://nodejs.org)
- Install [NPM](https://www.npmjs.com/package/npm) package manager
- Install required packages: 

```bash 
make build
```

## Instructions

To run the server you need simply to execute the following command:

```bash
make run
```

Or if you want to customize the port of the server (port 3001 for example), run:

```bash
PORT=3001 make run
```

Then you can access the data via the url: `http://localhost:${PORT}/api`

The api exposes the following endpoints:
```
GET /api/models
GET /api/models/{id}
GET /api/manufacturers
GET /api/manufacturers/{id}
GET /api/categories
GET /api/categories/{id}
```

The `image` property relates to an image for a `carModel`, and can be found in the `/api/images` directory.
