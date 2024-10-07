def print_env():
    print('''export POSTGRES_HOST="172.17.0.1"
export POSTGRES_PORT="5432"
export POSTGRES_PASSWORD="12345"
export POSTGRES_USER="gym"
export POSTGRES_DB="gym"
export PGDATA="/var/lib/postgresql/data/pgdata"
export POSTGRES_SSLMODE="disable"

export REDIS_HOST="172.17.0.1"
export REDIS_PORT="6379"
export REDIS_PASSWORD="12345"
export REDIS_DATABASES="16"

export GOLANG_HOST="172.17.0.1"
export GOLANG_PORT="8080"
export API_KEY="test_wuXPhl_ka-XYkEcO1mY27qN5RNYqOb1zVlc6Vogxlc4"
export SHOP_ID="395370"
          ''')
    

def main():
    print_env()


if __name__ == "__main__":
    main()