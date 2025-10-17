defmodule BackendIexWeb.PictureControllerTest do
  use BackendIexWeb.ConnCase

  import BackendIex.BackendFixtures
  alias BackendIex.Backend.Picture

  @create_attrs %{
    mode: "some mode",
    checksum: "some checksum",
    description: "some description",
    speed: "some speed",
    width: 42,
    ext: "some ext",
    camera: "some camera",
    aperture: "some aperture",
    iso: 42,
    focal_length: "some focal_length",
    lens: "some lens",
    flash: true,
    landscape: true,
    panoramic: true,
    height: 42,
    favourite: true,
    trigger_warning: true
  }
  @update_attrs %{
    mode: "some updated mode",
    checksum: "some updated checksum",
    description: "some updated description",
    speed: "some updated speed",
    width: 43,
    ext: "some updated ext",
    camera: "some updated camera",
    aperture: "some updated aperture",
    iso: 43,
    focal_length: "some updated focal_length",
    lens: "some updated lens",
    flash: false,
    landscape: false,
    panoramic: false,
    height: 43,
    favourite: false,
    trigger_warning: false
  }
  @invalid_attrs %{mode: nil, checksum: nil, description: nil, speed: nil, width: nil, ext: nil, camera: nil, aperture: nil, iso: nil, focal_length: nil, lens: nil, flash: nil, landscape: nil, panoramic: nil, height: nil, favourite: nil, trigger_warning: nil}

  setup %{conn: conn} do
    {:ok, conn: put_req_header(conn, "accept", "application/json")}
  end

  describe "index" do
    test "lists all pictures", %{conn: conn} do
      conn = get(conn, ~p"/api/pictures")
      assert json_response(conn, 200)["data"] == []
    end
  end

  describe "create picture" do
    test "renders picture when data is valid", %{conn: conn} do
      conn = post(conn, ~p"/api/pictures", picture: @create_attrs)
      assert %{"id" => id} = json_response(conn, 201)["data"]

      conn = get(conn, ~p"/api/pictures/#{id}")

      assert %{
               "id" => ^id,
               "aperture" => "some aperture",
               "camera" => "some camera",
               "checksum" => "some checksum",
               "description" => "some description",
               "ext" => "some ext",
               "favourite" => true,
               "flash" => true,
               "focal_length" => "some focal_length",
               "height" => 42,
               "iso" => 42,
               "landscape" => true,
               "lens" => "some lens",
               "mode" => "some mode",
               "panoramic" => true,
               "speed" => "some speed",
               "trigger_warning" => true,
               "width" => 42
             } = json_response(conn, 200)["data"]
    end

    test "renders errors when data is invalid", %{conn: conn} do
      conn = post(conn, ~p"/api/pictures", picture: @invalid_attrs)
      assert json_response(conn, 422)["errors"] != %{}
    end
  end

  describe "update picture" do
    setup [:create_picture]

    test "renders picture when data is valid", %{conn: conn, picture: %Picture{id: id} = picture} do
      conn = put(conn, ~p"/api/pictures/#{picture}", picture: @update_attrs)
      assert %{"id" => ^id} = json_response(conn, 200)["data"]

      conn = get(conn, ~p"/api/pictures/#{id}")

      assert %{
               "id" => ^id,
               "aperture" => "some updated aperture",
               "camera" => "some updated camera",
               "checksum" => "some updated checksum",
               "description" => "some updated description",
               "ext" => "some updated ext",
               "favourite" => false,
               "flash" => false,
               "focal_length" => "some updated focal_length",
               "height" => 43,
               "iso" => 43,
               "landscape" => false,
               "lens" => "some updated lens",
               "mode" => "some updated mode",
               "panoramic" => false,
               "speed" => "some updated speed",
               "trigger_warning" => false,
               "width" => 43
             } = json_response(conn, 200)["data"]
    end

    test "renders errors when data is invalid", %{conn: conn, picture: picture} do
      conn = put(conn, ~p"/api/pictures/#{picture}", picture: @invalid_attrs)
      assert json_response(conn, 422)["errors"] != %{}
    end
  end

  describe "delete picture" do
    setup [:create_picture]

    test "deletes chosen picture", %{conn: conn, picture: picture} do
      conn = delete(conn, ~p"/api/pictures/#{picture}")
      assert response(conn, 204)

      assert_error_sent 404, fn ->
        get(conn, ~p"/api/pictures/#{picture}")
      end
    end
  end

  defp create_picture(_) do
    picture = picture_fixture()

    %{picture: picture}
  end
end
