# ER Diagram

```mermaid
erDiagram
    USERS {
        string id PK
        string user_name
        string full_name
        string github_id
        string icon_url
    }

    CARDS {
        string id PK
        string user_id FK
        datetime created_at
    }

    COLLECTED_CARDS {
        string id PK
        string user_id FK
        string card_id FK
        datetime collected_at
    }

    IDENTICONS {
        string id PK
        string user_id FK
        string color
        json blocks_data
    }

    USERS ||--|| CARDS : creates
    USERS ||--o{ COLLECTED_CARDS : holds
    CARDS ||--o{ COLLECTED_CARDS : is_collected_in
    USERS ||--|| IDENTICONS : has
```
