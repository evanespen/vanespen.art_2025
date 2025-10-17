defmodule BackendIexWeb.Plug.Authenticate do
  @moduledoc """
  Authentication plug for protecting routes by verifying tokens.

  This plug:
  - Extracts the authorization token from the request header
  - Verifies the token using BackendIex.Utils.Security
  - Assigns current user data to the connection if token is valid
  - Returns 401 Unauthorized response if token is invalid or missing
  """

  import Plug.Conn
  require Bandit.Logger

  def init(opts) do
    opts
  end

  def call(conn, _opts) do
    with [token] <- get_req_header(conn, "authorization"),
         {:ok, data} <- BackendIex.Utils.Security.verify(token) do
      conn
      |> assign(:current_user, data.user)
    else
      _error ->
        conn
        |> Plug.Conn.send_resp(:unauthorized, "")
        |> halt()
    end
  end
end