defmodule BackendIex.Repo.Migrations.CreatePictures do
  use Ecto.Migration

  # TODO: Add :unique on checksum

  def change do
    create table(:pictures, primary_key: false) do
      add :id, :binary_id, primary_key: true

      add :timestamp, :integer
      add :ext, :string
      add :checksum, :string, unique: true
      add :camera, :string
      add :mode, :string
      add :aperture, :string
      add :iso, :integer
      add :speed, :string
      add :focal_length, :string
      add :lens, :string
      add :flash, :boolean
      add :landscape, :boolean
      add :panoramic, :boolean
      add :width, :integer
      add :height, :integer
      add :favourite, :boolean
      add :trigger_warning, :boolean
      add :description, :string

#      timestamps(type: :utc_datetime)
    end

    create unique_index(:pictures, [:checksum])
  end
end
