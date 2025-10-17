defmodule BackendIex.BackendTest do
  use BackendIex.DataCase

  alias BackendIex.Backend

  describe "pictures" do
    alias BackendIex.Backend.Picture

    import BackendIex.BackendFixtures

    @invalid_attrs %{}

    test "list_pictures/0 returns all pictures" do
      picture = picture_fixture()
      assert Backend.list_pictures() == [picture]
    end

    test "get_picture!/1 returns the picture with given id" do
      picture = picture_fixture()
      assert Backend.get_picture!(picture.id) == picture
    end

    test "create_picture/1 with valid data creates a picture" do
      valid_attrs = %{}

      assert {:ok, %Picture{} = picture} = Backend.create_picture(valid_attrs)
    end

    test "create_picture/1 with invalid data returns error changeset" do
      assert {:error, %Ecto.Changeset{}} = Backend.create_picture(@invalid_attrs)
    end

    test "update_picture/2 with valid data updates the picture" do
      picture = picture_fixture()
      update_attrs = %{}

      assert {:ok, %Picture{} = picture} = Backend.update_picture(picture, update_attrs)
    end

    test "update_picture/2 with invalid data returns error changeset" do
      picture = picture_fixture()
      assert {:error, %Ecto.Changeset{}} = Backend.update_picture(picture, @invalid_attrs)
      assert picture == Backend.get_picture!(picture.id)
    end

    test "delete_picture/1 deletes the picture" do
      picture = picture_fixture()
      assert {:ok, %Picture{}} = Backend.delete_picture(picture)
      assert_raise Ecto.NoResultsError, fn -> Backend.get_picture!(picture.id) end
    end

    test "change_picture/1 returns a picture changeset" do
      picture = picture_fixture()
      assert %Ecto.Changeset{} = Backend.change_picture(picture)
    end
  end

  describe "albums" do
    alias BackendIex.Backend.Album

    import BackendIex.BackendFixtures

    @invalid_attrs %{description: nil, title: nil, pictures: nil}

    test "list_albums/0 returns all albums" do
      album = album_fixture()
      assert Backend.list_albums() == [album]
    end

    test "get_album!/1 returns the album with given id" do
      album = album_fixture()
      assert Backend.get_album!(album.id) == album
    end

    test "create_album/1 with valid data creates a album" do
      valid_attrs = %{description: "some description", title: "some title", pictures: []}

      assert {:ok, %Album{} = album} = Backend.create_album(valid_attrs)
      assert album.description == "some description"
      assert album.title == "some title"
      assert album.pictures == []
    end

    test "create_album/1 with invalid data returns error changeset" do
      assert {:error, %Ecto.Changeset{}} = Backend.create_album(@invalid_attrs)
    end

    test "update_album/2 with valid data updates the album" do
      album = album_fixture()
      update_attrs = %{description: "some updated description", title: "some updated title", pictures: []}

      assert {:ok, %Album{} = album} = Backend.update_album(album, update_attrs)
      assert album.description == "some updated description"
      assert album.title == "some updated title"
      assert album.pictures == []
    end

    test "update_album/2 with invalid data returns error changeset" do
      album = album_fixture()
      assert {:error, %Ecto.Changeset{}} = Backend.update_album(album, @invalid_attrs)
      assert album == Backend.get_album!(album.id)
    end

    test "delete_album/1 deletes the album" do
      album = album_fixture()
      assert {:ok, %Album{}} = Backend.delete_album(album)
      assert_raise Ecto.NoResultsError, fn -> Backend.get_album!(album.id) end
    end

    test "change_album/1 returns a album changeset" do
      album = album_fixture()
      assert %Ecto.Changeset{} = Backend.change_album(album)
    end
  end
end
