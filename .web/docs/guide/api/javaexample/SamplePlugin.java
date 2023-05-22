package sample;

import build.buf.gen.minekube.connect.v1alpha1.ConnectEndpointRequest;
import build.buf.gen.minekube.connect.v1alpha1.ConnectEndpointResponse;
import build.buf.gen.minekube.connect.v1alpha1.ConnectServiceGrpc.ConnectServiceBlockingStub;
import build.buf.gen.minekube.connect.v1alpha1.ListEndpointsRequest;
import build.buf.gen.minekube.connect.v1alpha1.ListEndpointsResponse;
import com.minekube.connect.api.InstanceHolder;
import java.util.stream.Collectors;
import org.bukkit.Bukkit;
import org.bukkit.plugin.java.JavaPlugin;

public class SamplePlugin extends JavaPlugin {

    @Override
    public void onEnable() {
        ConnectServiceBlockingStub client = InstanceHolder.getClients().getConnectServiceBlockingStub();

        // List all endpoints
        ListEndpointsRequest listReq = ListEndpointsRequest.newBuilder().build();
        ListEndpointsResponse listRes = client.listEndpoints(listReq);
        getLogger().info(
                "First page of active and accessible Endpoints: " + listRes.getEndpointsList());

        // Move all online players to another Endpoint
        ConnectEndpointRequest connectReq = ConnectEndpointRequest.newBuilder()
                .setEndpoint("another-endpoint")
                .addAllPlayers(Bukkit.getOnlinePlayers().stream()
                        .map(p -> p.getUniqueId().toString())
                        .collect(Collectors.toList())
                ).build();
        getLogger().info("Moving players: " + connectReq.getPlayersCount());
        ConnectEndpointResponse connectRes = client.connectEndpoint(connectReq);
        getLogger().info(connectRes.toString());
    }
}
