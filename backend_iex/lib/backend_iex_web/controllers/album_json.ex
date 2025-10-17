defmodule BackendIexWeb.AlbumJSON do
  alias BackendIex.Backend.Album

  @doc """
  Renders a list of albums.
  """
  def index(%{albums: albums}) do
    for(album <- albums, do: data(album))
  end

  @doc """
  Renders a single album.
  """
  def show(%{album: album}) do
    data(album)
  end

  defp data(%Album{} = album) do
    %{
      id: album.id,
      title: album.title,
      description: album.description,
      pictures: album.pictures
    }
  end
end
