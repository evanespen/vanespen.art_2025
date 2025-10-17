defmodule BackendIex.BackendFixtures do
  @moduledoc """
  This module defines test helpers for creating
  entities via the `BackendIex.Backend` context.
  """

  @doc """
  Generate a picture.
  """
  def picture_fixture(attrs \\ %{}) do
    {:ok, picture} =
      attrs
      |> Enum.into(%{

      })
      |> BackendIex.Backend.create_picture()

    picture
  end

  @doc """
  Generate a album.
  """
  def album_fixture(attrs \\ %{}) do
    {:ok, album} =
      attrs
      |> Enum.into(%{
        description: "some description",
        pictures: [],
        title: "some title"
      })
      |> BackendIex.Backend.create_album()

    album
  end
end
