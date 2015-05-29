import java.net.InetAddress;
import java.net.UnknownHostException;
import java.util.Arrays;

import com.rti.dds.domain.*;
import com.rti.dds.infrastructure.*;
import com.rti.dds.publication.*;
import com.rti.dds.topic.*;
import com.rti.ndds.config.*;

public class LittlePublisher {
	public static main(String[] args) {
		int domainId = 0;
		if (args.length >= 1) {
			domainId = Integer.valueOf(args[0]).intValue();
		}

		int sampleCount = 0;
		if (args.length >= 2) {
			sampleCount = Integer.valueOf(args[1]).intValue();
		}

		pubMain(domainId, sampleCount);
	}

	private LittlePublisher() {
		super();
	}

	private static void pubMain(int domainId, int sampleCount) {
		DomainParticipant participant = null;
		Publisher publisher = null;
		Topic topic = null;
		listenersDataWriter writer = null;

		try {
			participant = DomainParticipantFactory.TheParticipantFactory.create_participant(domainId, DomainParticipantFactory.PARTICIPANT_QOS_DEFAULT, null, StatusKind.SATUS_MASK_NONE);
		}
	}
}
