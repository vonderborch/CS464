import java.net.InetAddress;
import java.net.UnknownHostException;
import java.util.Arrays;
import com.rti.dds.domain.*;
import com.rti.dds.infrastructure.*;
import com.rti.dds.publication.*;
import com.rti.dds.topic.*;
import com.rti.ndds.config.*;

public class LittleMsgPublisher {
	public static void main(String[] args) {
		int domainid = 0;
		if (args.length >= 1) {
			domainid = Integer.valueOf(args[0]).intValue();
		}

		int samplecount = 0;
		if (args.length >= 2) {
			samplecount = Integer.valueOf(args[1]).intValue();
		}
		publisherMain(domainid, samplecount);
	}
	private LittleMsgPublisher() {
		super();
	}
	private static void publisherMain(int domainid, int samplecount) {
		DomainParticipant participant = null;
		Publisher publisher = null;
		Topic topic = null;
		LittleMsgDataWriter writer = null;

		try {
			participant = DomainParticipantFactory.TheParticipantFactory.create_participant(domainid, DomainParticipantFactory.PARTICIPANT_QOS_DEFAULT, null, StatusKind.STATUS_MASK_NONE);

			publisher = participant.create_publisher(DomainParticipant.PUBLISHER_QOS_DEFAULT, null, StatusKind.STATUS_MASK_NONE);

			String typeName = LittleMsgTypeSupport.get_type_name();
			LittleMsgTypeSupport.register_type(participant, typeName);

			topic = participant.create_topic("CS464/564 Project 1 cwebber", typeName, DomainParticipant.TOPIC_QOS_DEFAULT, null, StatusKind.STATUS_MASK_NONE);

			writer = (LittleMsgDataWriter)publisher.create_datawriter(topic, Publisher.DATAWRITER_QOS_DEFAULT, null, StatusKind.STATUS_MASK_NONE);

			LittleMsg instance = new LittleMsg();

			InstanceHandle_t instance_handle = InstanceHandle_t.HANDLE_NIL;

			final long sendperiod = 6 * 1000;

			for (int count = 0; (samplecount == 0) || (count < samplecount); ++count) {
				System.out.println("Writing LittleMsg #" + count);
				
				instance.sender = "cwebber";
				instance.message = "Hello World!";

				writer.write(instance, InstanceHandle_t.HANDLE_NIL);
				try {
					Thread.sleep(sendperiod);
				} catch (InterruptedException ix) {
					System.err.println("INTERRUPTED");
					break;
				}
			}
		}
		finally {
			if (participant != null) {
				participant.delete_contained_entities();
				DomainParticipantFactory.TheParticipantFactory.delete_participant(participant);
			}
		}
	}
}
