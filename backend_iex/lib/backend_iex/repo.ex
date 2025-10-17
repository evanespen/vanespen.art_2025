defmodule BackendIex.Repo do
  use Ecto.Repo,
    otp_app: :backend_iex,
    adapter: Ecto.Adapters.SQLite3
end
