defmodule BackendIex.Repo.Migrations.CreateAlbums do
  use Ecto.Migration

  # TODO: Add :unique on title

  def change do
    create table(:albums, primary_key: false) do
      add :id, :binary_id, primary_key: true
      add :title, :string, unique: true
      add :description, :string
      add :pictures, {:array, :binary_id}

      # timestamps(type: :utc_datetime)
    end

    create unique_index(:albums, [:title])
  end
end
