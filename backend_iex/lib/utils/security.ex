defmodule BackendIex.Utils.Security do
  @moduledoc """
  Security utility module for handling token-based authentication.
  Provides functions for signing and verifying Phoenix tokens with a 14-day expiration period.

  The module offers two main functions:
  - sign/1: Creates a signed token from the provided data
  - verify/1: Verifies a token and returns the original data if valid
  """

  @token_age_secs 14 * 86_400

  def sign(data) do
    Phoenix.Token.sign(BackendIexWeb.Endpoint, Application.get_env(:backend_iex, :token_salt), data)
  end

  def verify(token) do
    case Phoenix.Token.verify(BackendIexWeb.Endpoint, Application.get_env(:backend_iex, :token_salt), token,
           max_age: @token_age_secs
         ) do
      {:ok, data} -> {:ok, data}
      _error -> {:error, :unauthenticated}
    end
  end
end