# This file is responsible for configuring your application
# and its dependencies with the aid of the Config module.
#
# This configuration file is loaded before any dependency and
# is restricted to this project.

# General application configuration
import Config

admin_username = System.get_env("ADMIN_USERNAME")
admin_password = System.get_env("ADMIN_PASSWORD")
token_salt = System.get_env("TOKEN_SALT")

if admin_username == nil or admin_password == nil or token_salt == nil do
  raise """
  Error: Following env vars are required:
  - ADMIN_USERNAME
  - ADMIN_PASSWORD
  - TOKEN_SALT
  """
end

config :backend_iex,
  ecto_repos: [BackendIex.Repo],
  generators: [timestamp_type: :utc_datetime, binary_id: true],
  storage: [root: "STORAGE", full: "STORAGE/FULL", half: "STORAGE/HALF", thumb: "STORAGE/THUMB"],
  admin_username: admin_username,
  admin_password: admin_password,
  token_salt: token_salt

config :elixir, :time_zone_database, Tzdata.TimeZoneDatabase

# Configures the endpoint
config :backend_iex, BackendIexWeb.Endpoint,
  url: [host: "localhost"],
  adapter: Bandit.PhoenixAdapter,
  render_errors: [
    formats: [json: BackendIexWeb.ErrorJSON],
    layout: false
  ],
  pubsub_server: BackendIex.PubSub,
  live_view: [signing_salt: "YSfVvxVV"]

# Configures the mailer
#
# By default it uses the "Local" adapter which stores the emails
# locally. You can see the emails in your browser, at "/dev/mailbox".
#
# For production it's recommended to configure a different adapter
# at the `config/runtime.exs`.
config :backend_iex, BackendIex.Mailer, adapter: Swoosh.Adapters.Local

# Configures Elixir's Logger
config :logger, :default_formatter,
  format: "$time $metadata[$level] $message\n",
  metadata: [:request_id]

# Use Jason for JSON parsing in Phoenix
config :phoenix, :json_library, Jason

# Import environment specific config. This must remain at the bottom
# of this file so it overrides the configuration defined above.
import_config "#{config_env()}.exs"
