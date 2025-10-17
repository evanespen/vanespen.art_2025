defmodule BackendIexWeb.SessionController do
  use BackendIexWeb, :controller

  action_fallback BackendIexWeb.FallbackController

  def get_token(conn, %{"username" => username, "password" => password}) do
    IO.puts username
    IO.puts password

    IO.puts Application.get_env(:backend_iex, :admin_username)
    IO.puts Application.get_env(:backend_iex, :admin_password)

    if username == Application.get_env(:backend_iex, :admin_username) && password == Application.get_env(:backend_iex, :admin_password) do
      token = BackendIex.Utils.Security.sign(%{user: username})
      conn = Plug.Conn.put_resp_header(conn, "Authorization", token)
      send_resp(conn, :accepted, "")
    else
      send_resp(conn, :unauthorized, "")
    end
  end

  def verify(conn, _params) do
    send_resp(conn, :ok, "")
  end
end