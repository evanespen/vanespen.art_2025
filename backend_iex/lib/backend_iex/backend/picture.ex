defmodule BackendIex.Backend.Picture do
  @moduledoc """
  Schema representing a picture entity in the system with its metadata and technical properties.

  This schema stores detailed information about pictures, including:
  - Basic file information (ID, timestamp, extension, checksum)
  - Camera settings (camera model, mode, aperture, ISO, speed, focal length, lens)
  - Image properties (dimensions, orientation flags)
  - User preferences (favourite, trigger warning, description)

  The schema works in conjunction with ImageProcessor utility for:
  - Storing processed image metadata
  - Managing multiple resolution versions
  - Handling EXIF data extraction
  - Supporting image file operations

  Fields:
  - id: UUID of the picture (binary_id)
  - timestamp: Unix timestamp of when the picture was taken
  - ext: File extension
  - checksum: MD5 hash of the image file
  - camera: Camera model name
  - mode: Camera exposure program
  - aperture: Aperture value (f-stop)
  - iso: ISO sensitivity
  - speed: Shutter speed
  - focal_length: Focal length in millimeters
  - lens: Lens model
  - flash: Whether flash was used
  - landscape: Whether image is in landscape orientation
  - panoramic: Whether image has panoramic aspect ratio
  - width: Image width in pixels
  - height: Image height in pixels
  - favourite: User favourite flag
  - trigger_warning: Content warning flag
  - description: User-provided description
  """

  use Ecto.Schema
  import Ecto.Changeset

  @primary_key {:id, :binary_id, autogenerate: false}
  @foreign_key_type :binary_id
  schema "pictures" do
    field :timestamp, :integer
    field :ext, :string
    field :checksum, :string
    field :camera, :string
    field :mode, :string
    field :aperture, :string
    field :iso, :integer
    field :speed, :string
    field :focal_length, :string
    field :lens, :string
    field :flash, :boolean
    field :landscape, :boolean
    field :panoramic, :boolean
    field :width, :integer
    field :height, :integer
    field :favourite, :boolean
    field :trigger_warning, :boolean
    field :description, :string
#    timestamps(type: :utc_datetime)
  end

  @doc false
  def changeset(picture, attrs) do
    picture
    |> cast(attrs, [:id, :timestamp, :ext, :checksum, :camera, :mode, :aperture, :iso, :speed, :focal_length, :lens, :flash, :landscape, :panoramic, :width, :height, :favourite, :trigger_warning, :description])
    |> validate_required([:id, :timestamp, :ext, :checksum, :camera, :mode, :aperture, :iso, :speed, :focal_length, :lens, :flash, :landscape, :panoramic, :width, :height, :favourite, :trigger_warning, :description])
    |> unique_constraint([:id, :checksum])
  end
end

