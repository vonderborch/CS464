import java.net.InetAddress;
import java.net.UnknownHostException;
import java.util.Arrays;
import com.rti.dds.domain.*;
import com.rti.dds.infrastructure.*;
import com.rti.dds.subscription.*;
import com.rti.dds.topic.*;
import com.rti.ndds.config.*;

public class LittleMsgSubscriber {
	public static void main(String[] args) {
		int domainid = 0;
		if (args.length >= 1) {
			domainid = Integer.valueOf(args[0]).intValue();
		}
		int samplecount = 0;
		if (args.length >= 2) {
			samplecount = Integer.valueOf(args[1]).intValue();
		}

		subscriberMain(domainid, samplecount);
	}

	private LittleMsgSubscriber() {
		super();
	}

	private static void subscriberMain(int domainid, int samplecount) {
		DomainParticipant participant = null;
		Subscriber subscriber = null;
		Topic topic = null;
		DataReaderListener listener = null;
		LittleMsgDataReader reader = null;

		try {
			participant = DomainParticipantFactory.TheParticipantFactory.create_participant(domainid, DomainParticipantFactory.PARTICIPANT_QOS_DEFAULT, null, StatusKind.STATUS_MASK_NONE);
			if (participant == null) {
				System.err.println("create_participant error\n");
				return;
			}

			subscriber = participant.create_subscriber(DomainParticipant.SUBSCRIBER_QOS_DEFAULT, null, StatusKind.STATUS_MASK_NONE);
			if (subscriber == null) {
				System.err.println("create_subscriber error\n");
				return;
			}

			String typename = LittleMsgTypeSupport.get_type_name();
			LittleMsgTypeSupport.register_type(participant, typename);

			topic = participant.create_topic("CS464/564 Project 1 cwebber", typename, DomainParticipant.TOPIC_QOS_DEFAULT, null, StatusKind.STATUS_MASK_NONE);
			if (topic == null) {
				System.err.println("create_topic error\n");
				return;
			}

			listener = new LittleMsgListener();

			reader = (LittleMsgDataReader)subscriber.create_datareader(topic, Subscriber.DATAREADER_QOS_DEFAULT, listener, StatusKind.STATUS_MASK_ALL);
			if (reader == null) {
				System.err.println("create_datareader error\n");
				return;
			}

			for (int count = 0; (samplecount == 0) || (count < samplecount); ++count)
			{
				System.out.println("LittleMsg subscriber sleeping for 6 seconds...");

				try {
					Thread.sleep(6 * 1000);
				} catch (InterruptedException ix) {
					System.err.println("INTERRUPTED");
					break;
				}
			}
		} finally {
			if (participant != null) {
				participant.delete_contained_entities();
				DomainParticipantFactory.TheParticipantFactory.delete_participant(participant);
			}
		}
	}

	private static class LittleMsgListener extends DataReaderAdapter {
		LittleMsgSeq _dataseq = new LittleMsgSeq();
		SampleInfoSeq _infoseq = new SampleInfoSeq();
		int counter = 0;

		public void on_data_available(DataReader reader) {
			LittleMsgDataReader LittleMsgReader = (LittleMsgDataReader)reader;

			try {
				LittleMsgReader.take(_dataseq, _infoseq, ResourceLimitsQosPolicy.LENGTH_UNLIMITED, SampleStateKind.ANY_SAMPLE_STATE, ViewStateKind.ANY_VIEW_STATE, InstanceStateKind.ANY_INSTANCE_STATE);

				for(int i = 0; i < _dataseq.size(); ++i) {
					SampleInfo info = (SampleInfo)_infoseq.get(i);

					if (info.valid_data) {
						counter++;
						System.out.println(((LittleMsg)_dataseq.get(i)).toString("LittleMsg Received #" + counter, 0));
					}
				}
			} catch (RETCODE_NO_DATA nodata) {
				// no data!
			} finally {
				LittleMsgReader.return_loan(_dataseq, _infoseq);
			}

		}
	}
}
