Laundre - Laundry Management System

RUN
go run main.go

Postman Test Guide
- Login
    POST {url}/login
  
        body
        {
            "username": "admin",
            "password": "admin"
        }

        response
        {
            "token": "{token}",
            "user": {
                "id": 1,
                "username": "admin",
                "role": "admin",
                "branchId": null
            }
        }

- Logout
    POST {url}/api/logout
  
        Authorization: Bearer Token = {token}

        response
        {
            "message": "Successfully logged out"
        }

- Create Branch
    POST {url}/api/admin/branches
  
        Authorization: Bearer Token = {token}

        body
        {
            "name": "Dagoma",
            "address": "Podomoro Apartemen",
            "phone": "081234567890"
        }

- Get All Branch
    GET {url}/api/admin/branches
  
        Authorization: Bearer Token = {token}

- Get Branch by ID
    GET {url}/api/admin/branches/{id}
  
        Authorization: Bearer Token = {token}

- Update Branch by ID
    PUT {url}/api/admin/branches
  
        Authorization: Bearer Token = {token}
        
        body
        {
            "name": "Dahoma",
            "address": "Podomoro Apartement",
            "phone": "081234567890"
        }

- Delete Branch by ID
    DELETE {url}/api/admin/branches/{id}
  
        Authorization: Bearer Token = {token}

- Create User
    POST {url}/api/admin/users
  
        Authorization: Bearer Token = {token}
        
        body
        {
            "username": "Andi",
            "password": "staff",
            "role": "staf",
            "branch_id": 1,
            "status": "active"
        }

- Get All Users
    GET {url}/api/admin/users
  
        Authorization: Bearer Token = {token}

- Get User by ID
    Get {url}/api/admin/users/{id}
  
        Authorization: Bearer Token = {token}
