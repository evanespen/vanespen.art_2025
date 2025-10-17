defmodule BackendIexWeb.PictureJSON do
  alias BackendIex.Backend.Picture

  @doc """
  Renders a list of pictures.
  """
  def index(%{pictures: pictures}) do
    for(picture <- pictures, do: data(picture))
  end

  @doc """
  Renders a single picture.
  """
  def show(%{picture: picture}) do
    data(picture)
  end

  defp data(%Picture{} = picture) do
    %{
      id: picture.id,
      ext: picture.ext,
      checksum: picture.checksum,
      camera: picture.camera,
      mode: picture.mode,
      aperture: picture.aperture,
      iso: picture.iso,
      speed: picture.speed,
      focal_length: picture.focal_length,
      lens: picture.lens,
      flash: picture.flash,
      landscape: picture.landscape,
      panoramic: picture.panoramic,
      width: picture.width,
      height: picture.height,
      favourite: picture.favourite,
      trigger_warning: picture.trigger_warning,
      description: picture.description
    }
  end
end
