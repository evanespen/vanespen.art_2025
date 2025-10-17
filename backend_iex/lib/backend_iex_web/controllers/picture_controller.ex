defmodule BackendIexWeb.PictureController do
  use BackendIexWeb, :controller

  alias BackendIex.Backend
  alias BackendIex.Backend.Picture
  alias BackendIex.Utils.ImageProcessor

  action_fallback BackendIexWeb.FallbackController

  def get_all(conn, _params) do
    pictures = Backend.list_pictures()
    render(conn, :index, pictures: pictures)
  end

  def get_one(conn, %{"id" => id}) do
    try do
      picture = Backend.get_picture!(id)
      render(conn, :show, picture: picture)
    rescue
      _e in Ecto.NoResultsError -> send_resp(conn, :not_found, "Picture not found")
    end
  end

  def create(conn, %{"files" => upload_params}) do
    # FIXME: Due to call order, the picture is saved even if already present in the database so it causes duplication
    try do
      Backend.create_picture(ImageProcessor.process_image(upload_params.path, upload_params.filename))
      send_resp(conn, :created, "Picture saved")
    rescue
      _e in Ecto.ConstraintError -> send_resp(conn, :conflict, "Picture already exists")
      _e in BadSavedFileError -> send_resp(conn, :internal_server_error, "The file has not been saved correctly while upload")
      _e in ExifExtractError -> send_resp(conn, :internal_server_error, "Unable to extract exif data")
      _e in DateExtractError -> send_resp(conn, :internal_server_error, "Unable to parse datetime")
    end
  end

  def delete(conn, %{"id" => id}) do
    picture = Backend.get_picture!(id)

    with {:ok, %Picture{}} <- Backend.delete_picture(picture) do
      ImageProcessor.delete_image(picture)
      send_resp(conn, :no_content, "")
    end
  end
end
