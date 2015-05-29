import java.net.InetAddress;
import java.net.UnknownHostException;
import java.util.Arrays;
import java.io.*;
import java.util.Date;
import java.text.DateFormat;
import java.text.SimpleDateFormat;
import java.util.Calendar;

public class TopUtil
{
	public float cpuuage;
	public float memusage;
	public int proccount;
	public String time;
	
	private float prev_cpu_idle;
	private float prev_cpu_total;
	
	public TopUtil()
	{
        cpuuage = 0;
        memusage = 0;
        proccount = 0;
	}
	
	public void Update()
	{
		// update time
		DateFormat dateFormat = new SimpleDateFormat("yyyy/MM/dd HH:mm:ss");
		Date date = new Date();
		time = dateFormat.format(date);
		
		BufferedReader procreader = null;
		try
		{
			// update cpu usage
			procreader = new BufferedReader(new InputStreamReader(new FileInputStream("/proc/stat")));
			String txtline = procreader.readLine();
			if (txtline == null)
			{
				throw new Exception("Can't find /proc/stat!");
			}
			else
			{
				String[] CPU = txtline.split("\\s+");
				float idle = Float.parseFloat(CPU[4]);
				float total = Float.parseFloat(CPU[1]) + Float.parseFloat(CPU[2]) + Float.parseFloat(CPU[3]) + Float.parseFloat(CPU[4]);
				float dt = total - prev_cpu_total;
				cpuuage = dt == 0 ? 0 : (1000 * (dt - (idle - prev_cpu_idle)) / dt + 5) / 10;
				
				prev_cpu_idle = idle;
				prev_cpu_total = total;
								
				// update proc count
				while (txtline.indexOf("procs_running") == -1 && txtline != null)
				{
					txtline = procreader.readLine();
				}
				proccount = -1;
				if (txtline != null)
				{
					String[] CNT = txtline.split("\\s+");
					proccount = Integer.parseInt(CNT[1]);
				}
			}
			
			// update mem usage
			procreader = new BufferedReader(new InputStreamReader(new FileInputStream("/proc/meminfo")));
			txtline = procreader.readLine();
			if (txtline == null)
			{
				throw new Exception("Can't find /proc/meminfo!");
			}
			else
			{
				String[] MEM = txtline.split("\\s+");
				float total = Float.parseFloat(MEM[1]);
				float free = total;
				
				txtline = procreader.readLine();
				if (txtline != null)
				{
					String[] MEMB = txtline.split("\\s+");
					free = Float.parseFloat(MEMB[1]);
				}
				
				float diff = total - free;
				memusage = diff / total;
			}
		}
		catch (Exception ex)
		{
		}
		finally
		{
			if (procreader != null)
			{
				try { procreader.close(); }
				catch (Exception e) { }
			}
		}
	}
}

