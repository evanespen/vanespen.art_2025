defmodule BackendIexWeb.AlbumControllerTest do
  use BackendIexWeb.ConnCase

  import BackendIex.BackendFixtures
  alias BackendIex.Backend.Album

  @create_attrs %{
    description: "some description",
    title: "some title",
    pictures: []
  }
  @update_attrs %{
    description: "some updated description",
    title: "some updated title",
    pictures: []
  }
  @invalid_attrs %{description: nil, title: nil, pictures: nil}

  setup %{conn: conn} do
    {:ok, conn: put_req_header(conn, "accept", "application/json")}
  end

  describe "index" do
    test "lists all albums", %{conn: conn} do
      conn = get(conn, ~p"/api/albums")
      assert json_response(conn, 200)["data"] == []
    end
  end

  describe "create album" do
    test "renders album when data is valid", %{conn: conn} do
      conn = post(conn, ~p"/api/albums", album: @create_attrs)
      assert %{"id" => id} = json_response(conn, 201)["data"]

      conn = get(conn, ~p"/api/albums/#{id}")

      assert %{
               "id" => ^id,
               "description" => "some description",
               "pictures" => [],
               "title" => "some title"
             } = json_response(conn, 200)["data"]
    end

    test "renders errors when data is invalid", %{conn: conn} do
      conn = post(conn, ~p"/api/albums", album: @invalid_attrs)
      assert json_response(conn, 422)["errors"] != %{}
    end
  end

  describe "update album" do
    setup [:create_album]

    test "renders album when data is valid", %{conn: conn, album: %Album{id: id} = album} do
      conn = put(conn, ~p"/api/albums/#{album}", album: @update_attrs)
      assert %{"id" => ^id} = json_response(conn, 200)["data"]

      conn = get(conn, ~p"/api/albums/#{id}")

      assert %{
               "id" => ^id,
               "description" => "some updated description",
               "pictures" => [],
               "title" => "some updated title"
             } = json_response(conn, 200)["data"]
    end

    test "renders errors when data is invalid", %{conn: conn, album: album} do
      conn = put(conn, ~p"/api/albums/#{album}", album: @invalid_attrs)
      assert json_response(conn, 422)["errors"] != %{}
    end
  end

  describe "delete album" do
    setup [:create_album]

    test "deletes chosen album", %{conn: conn, album: album} do
      conn = delete(conn, ~p"/api/albums/#{album}")
      assert response(conn, 204)

      assert_error_sent 404, fn ->
        get(conn, ~p"/api/albums/#{album}")
      end
    end
  end

  defp create_album(_) do
    album = album_fixture()

    %{album: album}
  end
end
