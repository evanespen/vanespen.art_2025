defmodule BackendIex.Backend.Album do
  @moduledoc """
  Schema representing an album entity in the system for organizing pictures.

  This schema manages collections of pictures by storing:
  - Album metadata (ID, title, description)
  - References to pictures through an array of UUIDs

  The schema works in conjunction with the Picture schema to:
  - Organize pictures into logical groupings
  - Maintain collection metadata
  - Store picture references

  Fields:
  - id: UUID of the album (binary_id)
  - title: Album title
  - description: Album description
  - pictures: Array of picture UUIDs referencing Picture entities
  """

  use Ecto.Schema
  import Ecto.Changeset

  @primary_key {:id, :binary_id, autogenerate: true}
  @foreign_key_type :binary_id
  schema "albums" do
    field :title, :string
    field :description, :string
    field :pictures, {:array, Ecto.UUID}

#    timestamps(type: :utc_datetime)
  end

  @doc false
  def changeset(album, attrs) do
    album
    |> cast(attrs, [:title, :description, :pictures])
    |> validate_required([:title, :description, :pictures])
  end
end
