defmodule BackendIexWeb.CDNController do
  @moduledoc """
  Controller responsible for serving image files through CDN endpoints.
  Provides functionality to serve images in different resolutions (full, half, thumb)
  based on the requested parameters. Images are served from configured storage paths
  and return appropriate HTTP responses based on image existence and validity.
  """
  
  use BackendIexWeb, :controller

  alias BackendIex.Backend

  action_fallback BackendIexWeb.FallbackController

  @doc """
  Serves an image file through CDN endpoint based on the requested resolution.

  ## Parameters
    * conn - Connection struct representing the HTTP request
    * id - String representing the unique identifier of the picture
    * res - String specifying the resolution version ("full", "half", or "thumb")

  ## Returns
    * Sends the requested image file with :found status if exists
    * Sends :not_found status if picture or file doesn't exist
    * Sends :bad_request status if resolution parameter is invalid

  ## Examples
      GET /cdn/123?res=full
      GET /cdn/456?res=thumb
  """
  def serve(conn, %{"id" => id, "res" => res}) do
    try do
      picture = Backend.get_picture!(id)

      picture_path = Path.join(case res do
        "full" -> Application.get_env(:backend_iex, :storage)[:full]
        "half" -> Application.get_env(:backend_iex, :storage)[:half]
        "thumb" -> Application.get_env(:backend_iex, :storage)[:thumb]
        _ -> send_resp(conn, :bad_request, "unknown res #{res}")
      end,
        "#{id}.#{picture.ext}")

      if File.exists?(picture_path) do
        send_file(conn, :found, picture_path)
      else
        send_resp(conn, :not_found, "")
      end

    rescue
      _e in Ecto.NoResultsError -> send_resp(conn, :not_found, "")
    end
  end
end