defmodule BackendIexWeb.Router do
  use BackendIexWeb, :router

  pipeline :api do
    plug :accepts, ["json"]
  end

  pipeline :authenticated do
    plug BackendIexWeb.Plug.Authenticate
  end

  scope "/api", BackendIexWeb do
    pipe_through :api

    post "/token", SessionController, :get_token

    get "/pictures/:id", PictureController, :get_one
    get "/pictures", PictureController, :get_all

    get "/albums/:title", AlbumController, :get_one
    get "/albums", AlbumController, :get_all
  end

  scope "/api/admin", BackendIexWeb do
    pipe_through [:api, :authenticated]

    get "/verify", SessionController, :verify

    post "/pictures", PictureController, :create
    delete "/pictures/:id", PictureController, :delete

    post "/albums", AlbumController, :create
    patch "/albums/:title", AlbumController, :update
  end

  scope "/cdn" do
    pipe_through :api
    get "/:id", BackendIexWeb.CDNController, :serve
  end

  # Enable LiveDashboard and Swoosh mailbox preview in development
  if Application.compile_env(:backend_iex, :dev_routes) do
    # If you want to use the LiveDashboard in production, you should put
    # it behind authentication and allow only admins to access it.
    # If your application does not have an admins-only section yet,
    # you can use Plug.BasicAuth to set up some basic authentication
    # as long as you are also using SSL (which you should anyway).
    import Phoenix.LiveDashboard.Router

    scope "/dev" do
      pipe_through [:fetch_session, :protect_from_forgery]

      live_dashboard "/dashboard", metrics: BackendIexWeb.Telemetry
      forward "/mailbox", Plug.Swoosh.MailboxPreview
    end
  end
end
