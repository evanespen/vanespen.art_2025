defmodule BadSavedFileError do
  defexception message: "the file has not been saved correctly while upload"
end

defmodule ExifExtractError do
  defexception message: "unable to extract exif data"
end

defmodule DateExtractError do
  defexception message: "unable to parse datetime"
end

defmodule BackendIex.Utils.ImageProcessor do
  @moduledoc """
  Image processing utility module for handling image file operations and metadata extraction.

  This module provides functionality for:
  - Processing uploaded images and generating multiple resolutions
  - Extracting and formatting EXIF metadata from images
  - Managing image storage in different directories
  - Converting technical camera settings to human-readable format
  - Deleting images and their associated files

  The module handles various image-related operations including:
  - File validation and storage
  - Image resizing (full, half, and thumbnail resolutions)
  - EXIF data extraction (camera settings, timestamps)
  - Checksum generation
  - Image dimension analysis

  Raises:
  - BadSavedFileError: When file upload or processing fails
  - ExifExtractError: When EXIF data extraction fails
  - DateExtractError: When date parsing from EXIF fails
  """

  @doc """
  Converts a shutter speed value to a human-readable string representation.

  Takes a numeric shutter speed value (in APEX format) and converts it to a formatted string.
  - Returns "1 sec" for 1-second exposure
  - Returns number of seconds for exposures longer than 1 second
  - Returns fraction format (1/x) for exposures shorter than 1 second

  ## Examples

      iex> shutter_speed_to_string(0)
      "1 sec"
      iex> shutter_speed_to_string(-1)
      "2"
      iex> shutter_speed_to_string(1)
      "1/2"

  """
  @spec shutter_speed_to_string(number()) :: String.t()
  def shutter_speed_to_string(shutter_speed_value) do
    shutter_speed_seconds = 1 / :math.pow(2, shutter_speed_value)
    denominator = round(1 / shutter_speed_seconds)

    cond do
      denominator == 1 ->
        "1 sec"
      denominator < 1 ->
        "#{round(shutter_speed_seconds)}"
      :true ->
        "1/#{denominator}"
    end
  end

  def process_image(image_path, orig_filename) do

    storage_config = Application.get_env(:backend_iex, :storage)

    unless File.exists?(image_path) do
      raise BadSavedFileError
    end

    uuid = Ecto.UUID.generate()

    extension = Path.extname(orig_filename)
                |> String.replace(".", "")
    filename = "#{uuid}.#{extension}"

    # Create root storage directory
    unless File.exists?(storage_config[:root]) do File.mkdir(storage_config[:root]) end

    # Move file to full res dir
    unless File.exists?(storage_config[:full]) do File.mkdir(storage_config[:full]) end
    File.copy(image_path, Path.join(storage_config[:full], filename))

    # Opens image for computations
    image = case Vix.Vips.Image.new_from_file(Path.join(storage_config[:full], filename)) do
      {:ok, image} -> image
      {:error, _reason} -> raise BadSavedFileError
      _ -> raise BadSavedFileError
    end

    height = Image.height(image)
    width = Image.width(image)

    # Create and save half res file
    unless File.exists?(storage_config[:half]) do File.mkdir(storage_config[:half]) end
    Vix.Vips.Operation.thumbnail!(Path.join(storage_config[:full], filename), trunc(width / 2))
    |> Vix.Vips.Image.write_to_file(Path.join(storage_config[:half], filename))

    # Create and save thumb res file
    unless File.exists?(storage_config[:thumb]) do File.mkdir(storage_config[:thumb]) end
    Vix.Vips.Operation.thumbnail!(Path.join(storage_config[:full], filename), trunc(width / 10))
    |> Vix.Vips.Image.write_to_file(Path.join(storage_config[:thumb], filename))

    # Computes MD5
    checksum = File.read!(Path.join(storage_config[:full], filename))
               |> :erlang.md5()
               |> Base.encode16(case: :lower)

    # Extracts exif data
    exif_data = case Image.exif(image) do
      {:ok, exif_data} -> exif_data
      {:error, _reason} -> raise ExifExtractError
      _ -> raise ExifExtractError
    end

    #    field :ext, :string
    #    field :checksum, :string
    #    field :camera, :string
    #    field :mode, :string
    #    field :aperture, :string
    #    field :iso, :integer
    #    field :speed, :string
    #    field :focal_length, :string
    #    field :lens, :string
    #    field :flash, :boolean
    #    field :landscape, :boolean
    #    field :panoramic, :boolean
    #    field :width, :integer
    #    field :height, :integer
    #    field :favourite, :boolean
    #    field :trigger_warning, :boolean
    #    field :description, :string

    timestamp = case DateTime.from_naive(exif_data[:exif][:datetime_original], "Europe/Paris") do
      {:ok, datetime_local} -> DateTime.to_unix(datetime_local, :second)
      {:error, _reason} -> raise DateExtractError
      _ -> raise DateExtractError
    end

    %{
      id: uuid,
      timestamp: timestamp,
      ext: extension,
      checksum: checksum,
      camera: exif_data[:model],
      mode: exif_data[:exif][:exposure_program],
      aperture: "f/#{exif_data[:exif][:aperture_value] |> Decimal.from_float() |> Decimal.round(1)}",
      iso: exif_data[:exif][:iso_speed_ratings],
      speed: shutter_speed_to_string(exif_data[:exif][:shutter_speed_value]),
      focal_length: "#{exif_data[:exif][:focal_length]}mm",
      lens: exif_data[:exif][:lens_model],
      flash: !String.contains?(exif_data[:exif][:flash], "Off"),
      landscape: width > height,
      panoramic: width / height > 1.5,
      width: width,
      height: height,
      favourite: false,
      trigger_warning: false,
      description: "none"
    }
  end

  def delete_image(picture) do
    filename = "#{picture.id}.#{picture.ext}"

    storage_config = Application.get_env(:backend_iex, :storage)

    File.rm!(Path.join(storage_config[:full], filename))
    File.rm!(Path.join(storage_config[:half], filename))
    File.rm!(Path.join(storage_config[:thumb], filename))
  end
end