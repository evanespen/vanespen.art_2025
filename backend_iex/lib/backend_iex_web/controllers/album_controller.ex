defmodule BackendIexWeb.AlbumController do
  use BackendIexWeb, :controller

  alias BackendIex.Backend
  alias BackendIex.Backend.Album

  action_fallback BackendIexWeb.FallbackController

  def get_all(conn, _params) do
    albums = Backend.list_albums()
    render(conn, :index, albums: albums)
  end

  def get_one(conn, %{"title" => title}) do
    try do
      album = Backend.get_album!(title)
      render(conn, :show, album: album)
    rescue
      _e in Ecto.NoResultsError -> send_resp(conn, :not_found, "Album not found")
    end
  end

  def create(conn, %{"title" => album_title, "description" => album_description}) do
    album_params = %{
      title: album_title,
      description: album_description,
      pictures: []
    }

    try do
      with {:ok, %Album{} = album} <- Backend.create_album(album_params) do
        conn
        |> put_status(:created)
        |> put_resp_header("location", ~p"/api/albums/#{album}")
        |> render(:show, album: album)
      end
    rescue
      _e in Ecto.ConstraintError -> send_resp(conn, :conflict, "Album with this title already exists")
    end
  end

  def update(conn, %{"title" => title, "pictures" => pictures}) do
    try do
      album = Backend.get_album!(title)
      with {:ok, %Album{} = album} <- Backend.update_album(album, %{pictures: pictures}) do
        render(conn, :show, album: album)
      end
    rescue
      _e in Ecto.NoResultsError -> send_resp(conn, :not_found, "Album not found")
    end
  end

  def delete(conn, %{"title" => title}) do
    album = Backend.get_album!(title)

    with {:ok, %Album{}} <- Backend.delete_album(album) do
      send_resp(conn, :no_content, "")
    end
  end
end
